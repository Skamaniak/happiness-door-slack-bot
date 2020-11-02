import React, {Component} from "react";
import Spinner from "react-spinkit";

class LoadingIndicator extends Component {

  render() {
    return (
      <>
        <div>Loading data about the meeting...</div>
        <Spinner name="three-bounce"/>
      </>
    )
  }
}

export default LoadingIndicator;
