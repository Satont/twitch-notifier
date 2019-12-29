import { User } from "../models/User";

declare module 'telegraf' {
  interface ContextMessageUpdate {
    public userDb: User
  } 
}