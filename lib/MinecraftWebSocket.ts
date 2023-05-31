import { WebSocket, WebSocketServer } from "ws";

// Imports for loading plugins
import { readdirSync } from "fs";
import path, { join } from "path";
import { fileURLToPath } from 'url';

// Bedrock server import
import { BedrockServer } from "./BedrockServer.js";

// General type imports
import { BedrockEvent, EventName } from "./events/Events.js";
import { Listener } from "./listeners/Listeners.js";


export class MinecraftWebSocket {
    // Parameters
    public wss: WebSocketServer;
    private eventListeners: any[] = [];
    private servers: any = {};

    constructor(WEBSOCKET_PORT: number) {
        // Create web socket server
        this.wss = new WebSocketServer({ port: WEBSOCKET_PORT }, () => {
            console.log(`MC BE Management Web Socket running on port ${WEBSOCKET_PORT}`);
        });

        // On connection handler
        this.wss.on('connection', this.onConnection.bind(this));
    }

    async onConnection(ws: WebSocket, req) {
        var url = req.url;
        var userID = url.slice(1);
        this.servers[userID] = new BedrockServer(userID, ws);
        console.log('Connected: ' + userID);

        // Subscribe to events
        for (var eventListener in this.eventListeners) {
            await this.servers[userID].subscribeToEvent(
                this.eventListeners[eventListener].eventName,
                this.eventListeners[eventListener].callback
            );
        }

        ws.on('close', () => {
            delete this.servers[userID]
            console.log('Disconnected: ' + userID)
        });
    }

    // Event listener
    async on(eventName: EventName, callback: (event: BedrockEvent) => void) {
        this.eventListeners.push({ eventName: eventName, callback: callback })
    }

    // Send command
    async sendCommand(server: string, command: string) {
        return await this.servers[server].sendCommand(command);
    }

    // Load listeners { eventName: EventName, callback: (event: BedrockEvent) => void) }
    async loadListeners(listeners: Listener[]) {
        for (var listener in listeners) {
            this.eventListeners.push(listeners[listener]);
        }
    }

    // Load plugins -- Dynamic imports don't work, so I'm benching this for now
    async loadPlugins(pluginsFolderName: string) {
        // Get path to plugins folder
        const __filename = fileURLToPath(import.meta.url);
        const __dirname = path.dirname(__filename);
        const pluginsDir: string = join(__dirname, "../" + pluginsFolderName);

        // Get plugins
        const plugins: string[] = readdirSync(pluginsDir);

        // Load listeners from plugins
        let pluginListeners: Listener[] = [];
        for (const plugin of plugins) {
            const pluginPath: string = join(pluginsDir, plugin);

            const { listeners } = await import(pluginPath);
            pluginListeners = pluginListeners.concat(listeners);

            console.log("Loaded " + pluginListeners.length + " listeners from " + plugin);
        }

        console.log("Loaded " + pluginListeners.length + " listeners from " + plugins.length + " plugins");

        // Load listeners into server
        await this.loadListeners(pluginListeners);
    }
}