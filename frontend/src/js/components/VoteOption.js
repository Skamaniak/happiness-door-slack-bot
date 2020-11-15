import React, {Component} from "react";
import PropTypes from "prop-types";
import styles from "../../css/modules/vote.option.module.css"
import {Row, Col, Button} from "react-bootstrap";

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

  render() {
    let {emojiUrl, optionText, voters, onVote} = this.props
    voters = voters || [];
    return (
      <>
        <Row className={styles.voteRow}>
          <Col md={{ span: 5, offset: 3}}>
            <img src={emojiUrl} alt={optionText} title={optionText}/>
            {optionText}
          </Col>
          <Col >
            <Button size="sm" variant="outline-dark" onClick={onVote}>Select</Button>
          </Col>
        </Row>
        <Row>
          <Col md={{ span: 5, offset: 3 }}>
            {this.renderVoterIcons(voters)}
            {this.renderVoteCount(voters)}
          </Col>
        </Row>
      </>
    );
  }
}

VoteOption.propTypes = {
  emojiUrl: PropTypes.string,
  optionText: PropTypes.string,
  voters: PropTypes.arrayOf(PropTypes.object),
  onVote: PropTypes.func
};

export default VoteOption;

