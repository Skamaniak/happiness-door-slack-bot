import React, {Component} from 'react';
import Spinner from 'react-bootstrap/Spinner';

class LoadingIndicator extends Component {

  render() {
    return (
      <div className="d-flex justify-content-center">
        <Spinner animation="border" role="status">
          <span className="sr-only">Loading...</span>
        </Spinner>
      </div>
    );
  }
}

export default LoadingIndicator;
