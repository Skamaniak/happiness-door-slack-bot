import React, {Component} from "react";
import App from "../components/App";
import UserStore from "../userStore";
import SlackUserBar from "../components/SlackUserBar";
import {Alert, Col, Container, Row} from "react-bootstrap";

export default class MainView extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userSet: UserStore.isUserSet()
    }
  }

  onUserChange() {
    this.setState({
      userSet: true
    })
  }

  renderUserSelection() {
    return (
      <>
        <SlackUserBar onUserChange={() => this.onUserChange()}/>
        <Row>
          <Col>
            <Alert variant={"info"}>
              You do not have your email address set yet. Email address is used to identify Slack user so we can do
              bi-directional state syncing. Please fill in your email address in the top right corner. We do not
              validate if it is really your email address so please be honest and do not try to impersonate
              other people.
            </Alert>
          </Col>
        </Row>
      </>
    );
  }

  renderApp() {
    return (<App/>);
  }

  render() {
    let mainComponent;
    if (!this.state.userSet) {
      mainComponent = this.renderUserSelection();
    } else {
      mainComponent = this.renderApp();
    }
    return (
      <Container>{mainComponent}</Container>
    );
  }

}