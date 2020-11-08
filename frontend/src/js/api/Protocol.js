import {getHappinessDoorId, getJsonCookie} from "../util/auth";

export const MessageType = {
  fromBackend: {
    happinessDoorData: 'HappinessDoorData'
  },
  toBackend: {
    voting: 'VoteAction'
  }

}

export const VotingAction = {
  happy: "VOTE_HAPPY",
  neutral: "VOTE_NEUTRAL",
  sad: "VOTE_SAD"
}

// TODO: this needs to be taken from OAUTH or cookie
export const createVoteMessage = (action) => {
  const user = getJsonCookie("happiness-door-user")

  return {
    "user": user,
    "actions": [
      {
        "action_id": action,
        "value": getHappinessDoorId()
      }
    ]
  }
}