import React, {Component} from 'react';

import {createVoteMessage, MessageType} from '../api/Protocol';
import {getAuthUrlParams} from '../util/auth';
import {config} from '../config';
import UserStore from '../userStore';
import Socket from '../api/Socket';
import LoadingIndicator from './LoadingIndicator';
import NoAccess from './NoAccess';
import HappinessDoor from './HappinessDoor';
import SlackUserBar from './SlackUserBar';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      connected: false,
      error: false,
    };
  }

  backendUrl() {
    let authParams = getAuthUrlParams();
    if (authParams) {
      authParams += '&u=' + UserStore.getUser();
    }
    return config.backendUrl + authParams;
  }

  registerHandlers(socket) {
    socket.on(MessageType.fromBackend.happinessDoorData, (data) => this.happinessDoorData(data));
  }

  componentDidMount() {
    this.connect();
  }

  connect() {
    let socket = new Socket(this.backendUrl());
    this.setState({socket});

    socket.on('connect', () => this.onConnect());
    socket.on('disconnect', (e) => this.onDisconnect(e));
    socket.on('error', (e) => this.onError(e));

    this.registerHandlers(socket);
  }

  disconnect() {
    if (this.state.connected) {
      this.state.socket.close();
    }
    this.setState({error: false});
  }

  reconnect() {
    this.disconnect();
    this.connect();
  }

  onConnect() {
    this.setState({connected: true});
  }

  onDisconnect(e) {
    this.setState({connected: false});

    if (e.code === 1006 && !this.state.error) {
      this.reconnect();
    }
  }

  onError() {
    this.setState({error: true});
  }

  onVote(action) {
    if (this.state.connected) {
      const message = createVoteMessage(action);
      this.state.socket.emit(MessageType.toBackend.voting, message);
    }
  }

  happinessDoorData(data) {
    this.setState({happinessDoor: data});
  }

  renderLoadingIndicator() {
    return (
      <LoadingIndicator/>
    );
  }

  renderAccessDenied() {
    return (
      <>
        <SlackUserBar onUserChange={() => this.reconnect()}/>
        <NoAccess/>
      </>
    );
  }

  render() {
    if (this.state.error) {
      return this.renderAccessDenied();
    }
    if (!this.state.connected || !this.state.happinessDoor) {
      return this.renderLoadingIndicator();
    }
    return (
      <>
        <SlackUserBar onUserChange={() => this.reconnect()}/>
        <HappinessDoor happinessDoor={this.state.happinessDoor}
          onVote={(a) => this.onVote(a)}
          onUserChange={() => this.reconnect()}/>
      </>
    );
  }
}

export default App;