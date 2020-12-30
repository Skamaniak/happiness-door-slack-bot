import {EventEmitter} from 'events';

// Socket class to construct and provide methods for WebSocket connections.
export default class Socket {
  constructor(url, eventEmitter = new EventEmitter()) {
    this.webSocket = new WebSocket(url);
    this.eventEmitter = eventEmitter;
    this.webSocket.onmessage = this.message.bind(this);
    this.webSocket.onopen = this.onOpen.bind(this);
    this.webSocket.onclose = this.onClose.bind(this);
    this.webSocket.onerror = this.onError.bind(this);
  }

  on(name, fn) {
    this.eventEmitter.on(name, fn);
  }

  off(name, fn) {
    this.eventEmitter.removeListener(name, fn);
  }

  emit(name, data) {
    const message = JSON.stringify({name, data});
    this.webSocket.send(message);
  }

  close() {
    this.webSocket.close(1000, 'Graceful disconnect');
  }

  onOpen() {
    this.eventEmitter.emit('connect');
  }

  onClose(e) {
    this.eventEmitter.emit('disconnect', e);
  }

  onError(e) {
    console.error('websocket error', e);
    this.eventEmitter.emit('error', e);
  }

  // message handles an incoming message and forwards it to an event listener.
  message(e) {
    try {
      const message = JSON.parse(e.data);
      this.eventEmitter.emit(message.name, message.data);
    } catch (err) {
      this.eventEmitter.emit('error', err);
      console.log(Date().toString() + ': ', err);
    }
  }
}