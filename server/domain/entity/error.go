package entity

type customError string

func (err customError) Error() string { // customError implements error interface
	return string(err)
}

const TransactionBeginError customError = "Could not start transaction"
const TransactionCommitError customError = "Could not commit transaction"

const UserNotFoundError customError = "Could not find user"
const UserConflictError customError = "Could not add user due to fields conflicting"
const UserConflictNotFoundError customError = "Could not find conflicing users"

const ForumNotFoundError customError = "Could not find forum"
const ForumConflictError customError = "Could not add forum due to fields conflicting"
const ForumConflictNotFoundError customError = "Could not find conflicting forum"

const ThreadNotFoundError customError = "Could not find thread"
const ThreadConflictError customError = "Could not add thread due to fields conflicting"
const ThreadConflictNotFoundError customError = "Could not find conflicting thread"

const PostNotFoundError customError = "Could not find post"
const ParentNotFoundError customError = "Could not find post's parent in post's thread"
