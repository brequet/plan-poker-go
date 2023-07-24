enum MessageType {
	JOIN_ROOM = 'join_room',
	USER_JOINED = 'user_joined',
	SUBMIT_ESTIMATE = 'submit_estimate',
	ESTIMATE_REVEALED = 'estimate_revealed',
	RESET_PLANNING = 'reset_planning'
}

interface MessageInterface {
	type: MessageType;
}

interface JoinRoomMessage extends MessageInterface {
	type: MessageType.JOIN_ROOM;
	payload: {
		roomName: string;
		userName: string;
	};
}

interface UserJoinedMessage {
	type: MessageType.USER_JOINED;
	payload: {
		userName: string;
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
	| ResetPlanningMessage;

export type {
	Message,
	JoinRoomMessage,
	UserJoinedMessage,
	SubmitEstimateMessage,
	EstimateRevealedMessage,
	ResetPlanningMessage
};
export { MessageType };
