import {EventEmitter} from 'events';

// Socket class to construct and provide methods for WebSocket connections.
export default class Socket {
  constructor(url, eventEmitter = new EventEmitter()) {
    try {
      this.webSocket = new WebSocket(url);
      this.eventEmitter = eventEmitter;
      this.webSocket.onmessage = this.message.bind(this);
      this.webSocket.onopen = this.open.bind(this);
      this.webSocket.onclose = this.close.bind(this);
      this.webSocket.onerror = this.error.bind(this);
    } catch (e) {
      console.error(e)
    }
  };

  // on adds a function as an event consumer/listener.
  on(name, fn) {
    this.eventEmitter.on(name, fn);
  };

  // off removes a function as an event consumer/listener.
  off(name, fn) {
    this.eventEmitter.removeListener(name, fn);
  };

  // open handles a connection to a websocket.
  open() {
    this.eventEmitter.emit('connect');
  };

  // close to handles a disconnection from a websocket.
  close() {
    this.eventEmitter.emit('disconnect');
  };

  // error handles an error on a websocket.
  error(e) {
    console.error("websocket error", e);
    this.eventEmitter.emit('error', e);
  }

  // emit sends a message on a websocket.
  emit(name, data) {
    const message = JSON.stringify({name, data});
    this.webSocket.send(message);
  }

  // message handles an incoming message and forwards it to an event listener.
  message(e) {
    try {
      const message = JSON.parse(e.data);
      this.eventEmitter.emit(message.name, message.data);
    } catch (err) {
      this.eventEmitter.emit('error', err);
      console.log(Date().toString() + ": ", err);
    }
  }
}