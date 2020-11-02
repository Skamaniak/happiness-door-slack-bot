import React, { Component } from "react";
import PropTypes from "prop-types";

class VoteOption extends Component {
  render() {
    let { emojiUrl, optionText } = this.props
    return (
      <div>
        <img src={emojiUrl}/>
        {optionText}
        <input type="button" value="Vote"/>
      </div>
    );
  }
}

VoteOption.propTypes = {
  emojiUrl: PropTypes.string,
  optionText: PropTypes.string
};

export default VoteOption;

