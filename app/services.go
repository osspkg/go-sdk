package app

import (
	"reflect"
	"sync/atomic"

	"github.com/deweppro/go-sdk/errors"
)

type (
	// ServiceInterface interface for services
	ServiceInterface interface {
		Up() error
		Down() error
	}
	//ServiceContextInterface interface for services with context
	ServiceContextInterface interface {
		Up(ctx Context) error
		Down() error
	}
)

var (
	srvType    = reflect.TypeOf(new(ServiceInterface)).Elem()
	srvTypeCtx = reflect.TypeOf(new(ServiceContextInterface)).Elem()
)

func asService(v reflect.Value) (ServiceInterface, bool) {
	if v.Type().AssignableTo(srvType) {
		return v.Interface().(ServiceInterface), true
	}
	return nil, false
}

func asServiceContext(v reflect.Value) (ServiceContextInterface, bool) {
	if v.Type().AssignableTo(srvTypeCtx) {
		return v.Interface().(ServiceContextInterface), true
	}
	return nil, false
}

func isService(v interface{}) bool {
	if _, ok := v.(ServiceInterface); ok {
		return true
	}
	if _, ok := v.(ServiceContextInterface); ok {
		return true
	}
	return false
}

/**********************************************************************************************************************/

const (
	statusUp   uint32 = 1
	statusDown uint32 = 0
)

type (
	_serv struct {
		tree   *treeItem
		status uint32
	}
	treeItem struct {
		Previous *treeItem
		Current  interface{}
		Next     *treeItem
	}
)

func newService() *_serv {
	return &_serv{
		tree:   nil,
		status: statusDown,
	}
}

// IsUp - mark that all services have started
func (s *_serv) IsUp() bool {
	return atomic.LoadUint32(&s.status) == statusUp
}

// Add - add new service by interface
func (s *_serv) Add(v interface{}) error {
	if s.IsUp() {
		return errDepRunning
	}

	if !isService(v) {
		return errors.Wrapf(errServiceUnknown, "service <%T>", v)
	}

	if s.tree == nil {
		s.tree = &treeItem{
			Previous: nil,
			Current:  v,
			Next:     nil,
		}
	} else {
		n := &treeItem{
			Previous: s.tree,
			Current:  v,
			Next:     nil,
		}
		n.Previous.Next = n
		s.tree = n
	}

	return nil
}

// Up - start all services
func (s *_serv) Up(ctx Context) error {
	if !atomic.CompareAndSwapUint32(&s.status, statusDown, statusUp) {
		return errDepRunning
	}
	if s.tree == nil {
		return nil
	}
	for s.tree.Previous != nil {
		s.tree = s.tree.Previous
	}
	for {
		if vv, ok := s.tree.Current.(ServiceContextInterface); ok {
			if err := vv.Up(ctx); err != nil {
				return err
			}
		} else if vv, ok := s.tree.Current.(ServiceInterface); ok {
			if err := vv.Up(); err != nil {
				return err
			}
		} else {
			return errors.Wrapf(errServiceUnknown, "service <%T>", s.tree.Current)
		}
		if s.tree.Next == nil {
			break
		}
		s.tree = s.tree.Next
	}

	return nil
}

// Down - stop all services
func (s *_serv) Down() (er error) {
	if !atomic.CompareAndSwapUint32(&s.status, statusUp, statusDown) {
		return errDepNotRunning
	}
	if s.tree == nil {
		return nil
	}
	for {
		if vv, ok := s.tree.Current.(ServiceContextInterface); ok {
			if err := vv.Down(); err != nil {
				er = errors.Wrap(er,
					errors.Wrapf(err, "down <%T> service error", s.tree.Current),
				)
			}
		} else if vv, ok := s.tree.Current.(ServiceInterface); ok {
			if err := vv.Down(); err != nil {
				er = errors.Wrap(er,
					errors.Wrapf(err, "down <%T> service error", s.tree.Current),
				)
			}
		} else {
			return errors.Wrapf(errServiceUnknown, "service <%T>", s.tree.Current)
		}
		if s.tree.Previous == nil {
			break
		}
		s.tree = s.tree.Previous
	}
	for s.tree.Next != nil {
		s.tree = s.tree.Next
	}
	return
}