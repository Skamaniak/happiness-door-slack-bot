import React, {Component} from "react";
import PropTypes from 'prop-types';

import styles from "../../css/modules/happiness.door.module.css"
import VoteOption from "./VoteOption";

class HappinessDoor extends Component {
  render() {
    let {meetingName} = this.props
    return (
      <>
        <div>How did you find the <span className={styles.meetingName}>{meetingName}</span> meeting?</div>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f642.png"}
                    optionText={"I'm happy"}/>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f610.png"}
                    optionText={"Neither good nor bad"}/>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f61e.png"}
                    optionText={"I did not like it"}/>
        <hr/>
      </>
    )
  }
}

HappinessDoor.propTypes = {
  meetingName: PropTypes.string,
};

export default HappinessDoor;

