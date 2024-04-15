import { FC, ReactNode, createContext, useState } from 'react';
import User from '../types/user';

export interface IUserContext {
  user: User;
  setUser: (user: User) => void;
}

export const UserContext = createContext<IUserContext | null>(null);

const UserProvider: FC<{children: ReactNode}> = ({children}) => {
    const [user, setUser] = useState<User>({} as User)

    return <UserContext.Provider value={{user, setUser}}>{children}</UserContext.Provider>
}

export default UserProvider