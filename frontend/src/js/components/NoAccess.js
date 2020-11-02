import React, {Component} from "react";

class NoAccess extends Component {

  render() {
    return (
      <>
        <div>You do not have an access to this page. Either the URL is wrong or access to this voting expired.</div>
      </>
    )
  }
}

export default NoAccess;
