export interface User {
    id: string;
    username: string;
}

export interface GameResult {
    username: string;
    wpm: number;
    accuracy: number;
}


export interface Player {
    username: string;
    progress: number;
}
