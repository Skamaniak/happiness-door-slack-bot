import React, {Component} from "react";
import Alert from 'react-bootstrap/Alert'

class NoAccess extends Component {

  render() {
    return (
      <Alert variant={"danger"}>
        <p>Server has some issues or you do not have access to this page. This can be caused by one of the following things</p>
        <ul>
          <li>The URL is wrong. Please make sure there is no typo.</li>
          <li>User emails you put in is not a recognised Slack user. Try to double-check it for typos.</li>
          <li>Backend threw an error. If you think the problem is not on your side, please contact the administrator or try to reload the page a tad later.</li>
        </ul>
      </Alert>
    )
  }
}

export default NoAccess;
