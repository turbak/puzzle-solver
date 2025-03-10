import { solveWithDisplay } from './algo';

let lastPostTime = 0;
const defaultThrottleInterval = 1000 / 10;

self.onmessage = async (event) => {
    const { month, day } = event.data;
    try {
        if (day < 1 || day > 31) {
            throw new Error('Invalid day');
        }

        let throttleInterval = event.data.throttleInterval;
        if (throttleInterval == null) {
            throttleInterval = defaultThrottleInterval;
        }

        const result = solveWithDisplay(month, day, (grid) => {
            const now = Date.now();
            if (now - lastPostTime >= throttleInterval) {
                lastPostTime = now;
                self.postMessage({ grid });
            }
        });
        self.postMessage({ result });
    } catch (error) {
        self.postMessage({ error: error });
    }
};