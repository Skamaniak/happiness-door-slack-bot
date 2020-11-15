import React, {Component} from "react";
import App from "../components/App";
import UserStore from "../userStore";
import SlackUserBar from "../components/SlackUserBar";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

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
            <p>You do not have your email address set yet. Email address is used to identify Slack user so we can do
              bi-directional state syncing. Fill in your email address first please.</p>
            <p>We do not validate if it is really your email address so please be honest and do not try to impersonate
              other people.</p>
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