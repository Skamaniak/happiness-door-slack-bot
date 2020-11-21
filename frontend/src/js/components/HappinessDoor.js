import React, {Component} from "react";
import PropTypes from 'prop-types';
import {VotingAction} from "../api/Protocol";

import styles from "../../css/modules/happiness.door.module.css"
import VoteOption from "./VoteOption";
import {Col, Row} from "react-bootstrap";

class HappinessDoor extends Component {
  render() {
    let {happinessDoor, onVote} = this.props

    return (
      <>
        <Row>
          <Col xl={{span: 6, offset: 3}} lg={{span: 12}} md={{span: 12}} sm={{span: 12}} xs={{span: 12}}>
            <div className={styles.title}>How did you find the <span
              className={styles.meetingName}>{happinessDoor.Name}</span> meeting?
            </div>
            <hr/>
          </Col>
        </Row>


        <VoteOption
          type={"happy"}
          optionText={"I'm happy"}
          voters={happinessDoor.HappyVoters}
          onVote={() => onVote(VotingAction.happy)}/>
        <VoteOption
          type={"neutral"}
          optionText={"Neither good nor bad"}
          voters={happinessDoor.NeutralVoters}
          onVote={() => onVote(VotingAction.neutral)}/>
        <VoteOption
          type={"sad"}
          optionText={"I did not like it"}
          voters={happinessDoor.SadVoters}
          onVote={() => onVote(VotingAction.sad)}/>
        <Col xl={{span: 6, offset: 3}} lg={{span: 12}} md={{span: 12}} sm={{span: 12}} xs={{span: 12}}>
          <hr/>
        </Col>
      </>
    )
  }
}

HappinessDoor.propTypes = {
  happinessDoor: PropTypes.object,
  onVote: PropTypes.func,
};

export default HappinessDoor;

