package main

type Stream struct {
	subscribers []*StreamSubscriber
}

func (s *Stream) IsActive() bool {
	return len(s.subscribers) > 0
}

func (s *Stream) Subscribe() *StreamSubscriber {
	ss := &StreamSubscriber{
		Flow:   make(chan string, 1),
		stream: s,
	}

	s.subscribers = append(s.subscribers, ss)
	return ss
}

func (s *Stream) NotifyAll(v string) {
	for _, ss := range s.subscribers {
		go ss.Notify(v)
	}
}

type StreamSubscriber struct {
	Flow   chan string
	stream *Stream
}

func (ss *StreamSubscriber) Notify(v string) {
	ss.Flow <- v
}

func (ss *StreamSubscriber) Unsubsribe() {
	list := ss.stream.subscribers
	close(ss.Flow)

	for i, x := range list {
		if x == ss {
			ss.stream.subscribers = append(list[0:i], list[i+1:]...)
			break
		}
	}
}
