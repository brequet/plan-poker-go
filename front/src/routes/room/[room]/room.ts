// src/stores/roomData.ts
import { writable, type Writable } from 'svelte/store';

export type User = {
	nickname: string;
	uuid?: string;
	estimate?: string;
};

export type CurrentUser = User & {
  isConnected: boolean;
}

export type Room = {
	code: string;
	name?: string;
	exist: boolean;
};

export const roomStore: Writable<Room> = writable();
export const currentUserStore: Writable<CurrentUser> = writable();
export const connectedUsersStore: Writable<User[]> = writable([]);

