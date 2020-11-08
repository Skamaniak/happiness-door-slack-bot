import React, {Component} from "react";

import {createVoteMessage, MessageType} from "../api/Protocol"

import Socket from '../api/Socket';
import LoadingIndicator from "./LoadingIndicator";
import NoAccess from "./NoAccess";
import HappinessDoor from "./HappinessDoor";
import {getAuthUrlParams} from "../util/auth";

export default class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      connected: false,
      error: false,
    }
  }

  backendUrl() {
    return 'ws://localhost:8081' + getAuthUrlParams() //TODO add to config
  }

  registerHandlers(socket) {
    socket.on(MessageType.fromBackend.happinessDoorData, (data) => this.happinessDoorData(data));
  }

  componentDidMount() {
    let socket = new Socket(this.backendUrl());
    this.setState({socket})

    socket.on('connect', () => this.onConnect());
    socket.on('disconnect', () => this.onDisconnect());
    socket.on('error', () => this.onError());
    this.registerHandlers(socket);
  }

  onConnect() {
    this.setState({connected: true});
  }

  onDisconnect() {
    this.setState({connected: false});
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
    this.setState({happinessDoor: data})
  }

  renderLoadingIndicator() {
    return (
      <LoadingIndicator/>
    )
  }

  renderAccessDenied() {
    return (
      <NoAccess/>
    )
  }

  render() {
    if (this.state.error) {
      return this.renderAccessDenied();
    }
    if (!this.state.connected || !this.state.happinessDoor) {
      return this.renderLoadingIndicator();
    }
    return (
      <HappinessDoor happinessDoor={this.state.happinessDoor} onVote={(a) => this.onVote(a)}/>
    )
  }
}