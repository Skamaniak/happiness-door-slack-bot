import React, {Component} from "react";
import PropTypes from 'prop-types';

import styles from "../../css/modules/happiness.door.module.css"
import VoteOption from "./VoteOption";
import {VotingAction} from "../api/Protocol";

class HappinessDoor extends Component {
  render() {
    let {happinessDoor, onVote} = this.props

    return (
      <>
        <div>How did you find the <span className={styles.meetingName}>{happinessDoor.Name}</span> meeting?</div>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f642.png"}
                    optionText={"I'm happy"}
                    voters={happinessDoor.HappyVoters}
                    onVote={() => onVote(VotingAction.happy)}/>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f610.png"}
                    optionText={"Neither good nor bad"}
                    voters={happinessDoor.NeutralVoters}
                    onVote={() => onVote(VotingAction.neutral)}/>
        <VoteOption emojiUrl={"https://a.slack-edge.com/production-standard-emoji-assets/10.2/google-medium/1f61e.png"}
                    optionText={"I did not like it"}
                    voters={happinessDoor.SadVoters}
                    onVote={() => onVote(VotingAction.sad)}/>
        <hr/>
      </>
    )
  }
}

HappinessDoor.propTypes = {
  happinessDoor: PropTypes.object,
  onVote: PropTypes.func
};

export default HappinessDoor;

