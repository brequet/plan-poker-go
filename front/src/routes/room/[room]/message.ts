enum MessageType {
	JOIN_ROOM = 'join_room',
	USER_JOINED = 'user_joined',
	SUBMIT_ESTIMATE = 'submit_estimate',
	ESTIMATE_REVEALED = 'estimate_revealed',
	RESET_PLANNING = 'reset_planning',
	USER_DISCONNECTED = 'user_disconnected',
	CONFIRM_CONNECTION = 'confirm_connection',
}

interface MessageInterface {
	type: MessageType;
}

interface UserDTO {
	userName: string;
	uuid: string;
}

interface UserDisconnectedMessage {
	type: MessageType.USER_DISCONNECTED;
	payload: {
		user: UserDTO;
	};
}
interface ConfirmConnectionMessage {
	type: MessageType.CONFIRM_CONNECTION;
	payload: {
		user: UserDTO;
		ConnectedUsers: UserDTO[]
	};
}

interface JoinRoomMessage extends MessageInterface {
	type: MessageType.JOIN_ROOM;
	payload: {
		roomCode: string;
		nickname: string;
	};
}

interface UserJoinedMessage {
	type: MessageType.USER_JOINED;
	payload: {
		user: UserDTO;
	};
}

interface SubmitEstimateMessage {
	type: MessageType.SUBMIT_ESTIMATE;
	payload: {
		taskId: string;
		estimate: number;
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
	| ConfirmConnectionMessage;

export type {
	Message,
	JoinRoomMessage,
	UserJoinedMessage,
	SubmitEstimateMessage,
	EstimateRevealedMessage,
	ResetPlanningMessage,
	UserDisconnectedMessage,
	ConfirmConnectionMessage
};
export { MessageType };
