package cudnn

/* Generated by gencudnn. DO NOT EDIT */

// #include <cudnn_v7.h>
import "C"
import "runtime"

// Dropout is a representation of cudnnDropoutDescriptor_t.
//
// The usecase of Dropout is quite different from the rest of the APIs in this library. There is a two stage API:
//		drop := NewDropout(...)
//		drop.Use(ctx, states....)
//
// This is because the Dropout is largely tied to run-time. An additional `.IsReady` method is added to indicate if the dropout state is ready to be used
//
// However, if your runtime is known ahead of time, the `NewDropoutWithContext` creation function can be used.
type Dropout struct {
	internal C.cudnnDropoutDescriptor_t

	handle           *Context
	dropout          float32
	states           Memory
	stateSizeInBytes uintptr
	seed             uint64
}

// NewDropout creates a Dropout descriptor. It is not usable by default because some additional stateful information needs to be passed in
func NewDropout(dropout float64) (retVal *Dropout, err error) {
	var internal C.cudnnDropoutDescriptor_t
	if err := result(C.cudnnCreateDropoutDescriptor(&internal)); err != nil {
		return nil, err
	}
	runtime.SetFinalizer(retVal, destroyDropout)
	retVal = &Dropout{
		internal: internal,
		dropout:  float32(dropout),
	}
	return retVal, nil
}

// NewDropout creates a new Dropout with the given context (handle, states, etc)
func NewDropoutWithContext(dropout float64, handle *Context, states Memory, stateSizeInBytes uintptr, seed uint64) (retVal *Dropout, err error) {
	retVal = NewDropoutDescriptor(dropout)
	err = retVal.Use(handle, states, stateSizeInBytes, seed)
	return
}

// Use is the second stage of the two-stage API.
func (d *Dropout) Use(ctx *Context, states Memory, stateSizeInBytes uintptr, seed uint64) error {
	d.handle = ctx
	d.states = states
	d.stateSizeInBytes = stateSizeInBytes
	d.seed = seed
	return result(C.cudnnSetDropoutDescriptor(internal, handle.internal, C.float(d.dropout), d.states.Pointer(), C.size_t(d.stateSizeInBytes), C.ulonglong(d.seed)))
}

// IsReady indicates if the dropout operator is ready to be used
func (d *Dropout) IsReady() bool {
	return d.handle != nil && d.states != nil && d.stateSizeInBytes != 0
}

// Reset resets the state to be not ready. It does NOT reset the dropout ratio.
func (d *Dropout) Reset() {
	d.handle = nil
	d.states = nil
	d.stateSizeInBytes = 0
	d.seed = 0
}

// Handle returns the internal handle.
func (d *Dropout) Handle() *Context { return d.handle }

// Dropout returns the internal dropout ratio.
func (d *Dropout) Dropout() float32 { return d.dropout }

// StateSizeInBytes returns the internal stateSizeInBytes.
func (d *Dropout) StateSizeInBytes() uintptr { return d.stateSizeInBytes }

// Seed returns the internal seed.
func (d *Dropout) Seed() uint64 { return d.seed }

func destroyDropout(obj *Dropout) { C.cudnnDestroyDropoutDescriptor(obj.internal) }
