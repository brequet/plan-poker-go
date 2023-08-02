enum MessageType {
	JOIN_ROOM = 'join_room',
	USER_JOINED = 'user_joined',
	SUBMIT_ESTIMATE = 'submit_estimate',
	ESTIMATE_REVEALED = 'estimate_revealed',
	RESET_PLANNING = 'reset_planning',
	USER_DISCONNECTED = 'user_disconnected',
	CONFIRM_CONNECTION = 'confirm_connection',
	CONFIRM_ESTIMATE_SUBMISSION = 'confirm_estimate_submission',
	USER_SUBMITTED_ESTIMATE = 'user_submitted_estimate'
}

interface MessageInterface {
	type: MessageType;
}

interface UserDTO {
	userName: string;
	uuid: string;
	estimate: string;
}

interface RoomDTO {
	roomCode: string;
	roomName: string;
}

// sended

interface JoinRoomMessage extends MessageInterface {
	type: MessageType.JOIN_ROOM;
	payload: {
		roomCode: string;
		nickname: string;
	};
}

interface SubmitEstimateMessage {
	type: MessageType.SUBMIT_ESTIMATE;
	payload: {
		estimate: string;
	};
}

// received

interface ConfirmConnectionMessage {
	type: MessageType.CONFIRM_CONNECTION;
	payload: {
		user: UserDTO;
		connectedUsers: UserDTO[];
		room: RoomDTO;
	};
}

interface UserDisconnectedMessage {
	type: MessageType.USER_DISCONNECTED;
	payload: {
		user: UserDTO;
	};
}

interface UserJoinedMessage {
	type: MessageType.USER_JOINED;
	payload: {
		user: UserDTO;
	};
}

interface ConfirmEstimateSubmittedMessage {
	type: MessageType.CONFIRM_ESTIMATE_SUBMISSION;
	payload: {
		estimate: string;
	};
}

interface UserSubmittedEstimateMessage {
	type: MessageType.USER_SUBMITTED_ESTIMATE;
	payload: {
		estimate: string;
		user: UserDTO;
	};
}

interface EstimateRevealedMessage {
	type: MessageType.ESTIMATE_REVEALED;
	payload: {
		estimates: Record<string, number>;
	};
}

interface ResetPlanningMessage {
	type: MessageType.RESET_PLANNING;
	payload: null;
}

// Define a union type for all possible message types
type Message =
	| JoinRoomMessage
	| UserJoinedMessage
	| SubmitEstimateMessage
	| EstimateRevealedMessage
	| ResetPlanningMessage
	| UserDisconnectedMessage
	| ConfirmConnectionMessage
	| ConfirmEstimateSubmittedMessage
	| UserSubmittedEstimateMessage;

export { MessageType };
export type { JoinRoomMessage, Message, SubmitEstimateMessage };

