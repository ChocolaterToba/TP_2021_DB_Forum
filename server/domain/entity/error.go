package entity

type customError string

func (err customError) Error() string { // customError implements error interface
	return string(err)
}

const TransactionBeginError customError = "Could not start transaction"
const TransactionCommitError customError = "Could not commit transaction"

const QueryParseError customError = "Could not parse query parameters"

const UnsupportedSortingModeError customError = "Passed sorting mode is incorrect or isn't supported yet"
const UnsupportedRelatedObjectError customError = "Passed related object is incorrect - can only be user, thread or forum"

const UserNotFoundError customError = "Could not find user"
const UserConflictError customError = "Could not add user due to fields conflicting"
const UserConflictNotFoundError customError = "Could not find conflicing users"

const ForumNotFoundError customError = "Could not find forum"
const ForumConflictError customError = "Could not add forum due to fields conflicting"
const ForumConflictNotFoundError customError = "Could not find conflicting forum"

const ThreadNotFoundError customError = "Could not find thread"
const ThreadConflictError customError = "Could not add thread due to fields conflicting"
const ThreadConflictNotFoundError customError = "Could not find conflicting thread"
const IncorrectVoteAmountError customError = "Vote values other than -1 and +1 are unacceptable"
const VoteNotFoundError customError = "Vote could not be found"
const VoteAlreadyExistsError customError = "Vote already exists"

const PostNotFoundError customError = "Could not find post"
const ParentNotFoundError customError = "Could not find post's parent"
const ParentInAnotherThreadError customError = "Post's parent is in different thread"

const InvalidIncrementValueError customError = "Cannot increment by non-positive value"
