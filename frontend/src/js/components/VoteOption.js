import React, {Component} from "react";
import PropTypes from "prop-types";
import styles from "../../css/modules/vote.option.module.css"

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
      <div>
        <div>
          <img src={emojiUrl} alt={optionText} title={optionText}/>
          {optionText}
          <input type="button" value="Vote" onClick={onVote}/>
        </div>

        {this.renderVoterIcons(voters)}
        {this.renderVoteCount(voters)}
      </div>
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

