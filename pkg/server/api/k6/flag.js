import { StringWithDefault } from './util.js'

// console.log(`-e BASE_URL=http://localhost:8099`)
export const baseURL = StringWithDefault(`http://localhost:8099`)(__ENV.BASE_URL)
