import React, {Component} from "react";
import PropTypes from "prop-types";
import styles from "../../css/modules/vote.option.module.css"
import {Row, Col, Button} from "react-bootstrap";
import HappyEmoji from "../../img/emoji_happy.png"
import NeutralEmoji from "../../img/emoji_neutral.png"
import SadEmoji from "../../img/emoji_sad.png"

class VoteOption extends Component {

  renderVoterIcons(voters) {
    return (
      <>
        {voters.map(voter => (
          <img key={voter.Id}
               className={styles.voterPicture}
               src={voter.ProfilePicture}
               alt={voter.Name}
               title={voter.Name}/>
        ))}
      </>
    )
  }

  renderVoteCount(voters) {
    if (voters.length > 0) {
      return (
        <>
          <span> - {voters.length} votes</span>
        </>
      )
    }
  }

  getEmoji(type) {
    if (type === "happy") {
      return HappyEmoji;
    } else if (type === "neutral") {
      return NeutralEmoji;
    }
    return SadEmoji;
  }

  render() {
    let {type, optionText, voters, onVote} = this.props
    voters = voters || [];

    return (
      <>
        <Row className={styles.voteRow}>
          <Col xl={{ span: 5, offset: 3}} lg={{ span: 10}} md={{ span: 10}} sm={{ span: 10 }} xs={{ span: 10 }}>
            <img src={this.getEmoji(type)} className={styles.emojiPicture} alt={optionText} title={optionText}/>
            &nbsp;
            <span>{optionText}</span>
          </Col>
          <Col >
            <Button size="sm" variant="outline-dark" onClick={onVote}>Select</Button>
          </Col>
        </Row>
        <Row>
          <Col xl={{ span: 5, offset: 3}} lg={{ span: 10}} md={{ span: 10}} sm={{ span: 10 }} xs={{ span: 10 }}>
            {this.renderVoterIcons(voters)}
            {this.renderVoteCount(voters)}
          </Col>
        </Row>
      </>
    );
  }
}

VoteOption.propTypes = {
  type: PropTypes.string,
  optionText: PropTypes.string,
  voters: PropTypes.arrayOf(PropTypes.object),
  onVote: PropTypes.func
};

export default VoteOption;

