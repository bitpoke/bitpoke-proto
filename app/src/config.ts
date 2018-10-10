declare global {
    interface Window {
        env: any
    }
}

const env = window.env
    || process.env
    || {}

export default env
