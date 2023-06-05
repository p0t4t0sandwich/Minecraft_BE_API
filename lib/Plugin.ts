// Plugin Class for extending the core functionality of the framework

import { MinecraftWebSocket } from "./MinecraftWebSocket.js";
import { Listener } from "./listeners/Listeners.js";

export class Plugin {
    // Properties
    readonly name: string;
    readonly description: string;
    readonly version: string;
    readonly author: string;

    private listeners: Listener[] = [];
    public mwss: MinecraftWebSocket;

    // Constructor
    constructor(name: string, description: string, version: string, author: string) {
        this.name = name;
        this.description = description;
        this.version = version;
        this.author = author;
    }

    // Methods

    // Get Listeners
    getListeners(): Listener[] {
        return this.listeners;
    }

    // Set Listeners
    setListeners(listeners: Listener[]) {
        this.listeners = listeners;
    }

    // Add Listener
    addListener(listener: Listener) {
        this.listeners.push(listener);
    }

    // Start
    async start(mwss: MinecraftWebSocket) {
        this.mwss = mwss;
        console.log("Plugin started");
    }
}