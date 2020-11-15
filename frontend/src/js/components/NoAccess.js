import React, {Component} from "react";

class NoAccess extends Component {

  render() {
    return (
      <>
        <p>You do not have an access to this page. This can be caused by one of the following things</p>
        <ul>
          <li>The URL is wrong</li>
          <li>Access to this happiness door voting has already expired</li>
          <li>User you put in is not a recognised Slack user</li>
        </ul>
      </>
    )
  }
}

export default NoAccess;
