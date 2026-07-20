package analytics

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Registry keeps info as per UTC times.
//
// MonthCurrent and MonthPrevious are swapped lock-free via atomic.Pointer.
// Ingestion and queries must go through CurrentMonth()/PreviousMonth() —
// never hold onto the returned pointer across a rollover boundary, and
// never store it; always re-call the accessor for each logical operation.
type Registry struct {
	monthCurrent  atomic.Pointer[MonthActive]
	monthPrevious atomic.Pointer[MonthActive]

	History [7]MonthArchived

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CalendarMonthCurrentNumber int8 // no need for year as we keep only 7 months.

	// rolloverMu serializes concurrent Rollover callers only. Ingestion
	// and queries never touch it and are never blocked by it.
	rolloverMu sync.Mutex
}

// NewRegistry returns a Registry ready for ingestion. The zero value of
// Registry is NOT usable — monthCurrent/monthPrevious must be seeded with
// real buffers before CurrentMonth()/PreviousMonth() can be called.
func NewRegistry(dstTimestamps ...int64) *Registry {
	r := Registry{}

	if len(dstTimestamps) == 2 {
		r.TimestampDSTSpring = dstTimestamps[0]
		r.TimestampDSTWinter = dstTimestamps[1]
	}

	r.monthCurrent.Store(&MonthActive{})
	r.monthPrevious.Store(&MonthActive{})

	return &r
}

// CurrentMonth returns the active buffer for the current month.
// Callers must not cache the returned pointer beyond a single logical
// operation — always re-call CurrentMonth() for the next one.
func (r *Registry) GetCurrentMonth() *MonthActive {
	return r.monthCurrent.Load()
}

// PreviousMonth returns the active buffer for the previous month.
// Same caching rule as CurrentMonth.
func (r *Registry) GetPreviousMonth() *MonthActive {
	return r.monthPrevious.Load()
}

// Rollover advances the registry by one month:
//  1. archives the current "previous" buffer into History[0] (shifting
//     the rest of the window down, dropping the oldest month)
//  2. deep-copies the (about to be retired) current buffer into a fresh
//     "previous" buffer
//  3. atomically swaps in a brand-new, zeroed current buffer
//
// Rollover is safe to call concurrently with ingestion: the swap is a
// single atomic pointer store, so CurrentMonth()/PreviousMonth() never
// observe a partially-reset buffer. It is NOT safe to call Rollover
// concurrently with itself — that is what rolloverMu is for.
//
// See the package doc for the one accepted race: in-flight increments
// that read a *MonthActive just before the swap and write to it just
// after are lost, since that buffer is never merged into history.
func (r *Registry) Rollover() {
	r.rolloverMu.Lock()
	defer r.rolloverMu.Unlock()

	retiredCurrent := r.monthCurrent.Load()
	retiredPrevious := r.monthPrevious.Load()

	// 1. Shift history window, then archive the retiring "previous" into
	//    History[0]. Built off to the side; History itself isn't touched
	//    until this is fully assembled.
	var newHistory [7]MonthArchived

	for ix := range 6 {
		newHistory[ix+1] = r.History[ix]
	}

	for day := range 31 {
		newHistory[0][day] = retiredPrevious[day].AsDayArchived()
	}

	r.History = newHistory

	// 2. Deep-copy the retiring current into a fresh previous buffer.
	newPrevious := &MonthActive{}
	retiredCurrent.DeepCopyInto(newPrevious)

	// 3. Publish. Both swaps are independent atomic stores; a reader can
	// briefly observe new current + old previous, but never a torn buffer.
	r.monthPrevious.Store(newPrevious)
	r.monthCurrent.Store(&MonthActive{})

	r.CalendarMonthCurrentNumber++
	if r.CalendarMonthCurrentNumber == 13 {
		r.CalendarMonthCurrentNumber = 1
	}
}

// Snapshot writes a textual snapshot of the registry into w.
func (r *Registry) Snapshot(w io.Writer) error {
	currentMonth := r.GetCurrentMonth()

	for day := range 31 {
		for hour := range 24 {
			m := &currentMonth[day][hour]

			if m.RecordsPerPeriod.Load() == 0 {
				continue
			}

			if _, errPrintCurrent := fmt.Fprintf(
				w,

				"current day:'%02d' hour:'%02d' records:'%d'\n",
				day,
				hour,
				m.RecordsPerPeriod.Load(),
			); errPrintCurrent != nil {
				return errPrintCurrent
			}
		}
	}

	previousMonth := r.GetPreviousMonth()

	for day := range 31 {
		for hour := range 24 {
			m := &previousMonth[day][hour]
			if m.RecordsPerPeriod.Load() == 0 {
				continue
			}

			if _, errPrintPrevious := fmt.Fprintf(
				w,

				"previous day: %02d hour: %02d records:%d\n",
				day,
				hour,
				m.RecordsPerPeriod.Load(),
			); errPrintPrevious != nil {
				return errPrintPrevious
			}
		}
	}

	for h := range 7 {
		for day := range 31 {
			for hour := range 24 {
				m := &r.History[h][day][hour]
				if m.RecordsPerPeriod == 0 {
					continue
				}

				if _, errPrintHistory := fmt.Fprintf(
					w,

					"history[%d] day: %02d hour: %02d records:%d\n",
					h,
					day,
					hour,
					m.RecordsPerPeriod,
				); errPrintHistory != nil {
					return errPrintHistory
				}
			}
		}
	}

	return nil
}
