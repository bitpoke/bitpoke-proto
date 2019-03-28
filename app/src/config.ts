import { isEmpty } from 'lodash'

declare global {
    interface Window {
        env: any
    }
}

const env = !isEmpty(window.env)
    ? window.env
    : process.env || {}

export default env
