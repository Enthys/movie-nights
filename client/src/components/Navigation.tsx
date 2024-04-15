import { useContext } from "react"
import { isLoggedIn } from "../utils"
import { IUserContext, UserContext } from "../context/userContext";
import User from "../types/user";
import { Link, redirect } from "react-router-dom";

class UserNotLoggedError extends Error {
  constructor() {
    super("user is not logged in");
  }
}

let loginAttemptMade = false;

async function getCurrentUser(): Promise<User> {
  return new Promise<User>((resolve) => {
    const resp = fetch("/api/profile")
      .then((resp) => {
        if (resp.status === 401) {
          return new UserNotLoggedError()
        }
        return resp.json()
      });

    resp.then((data) => {
      console.log(data);
      resolve(data.user);
    });
  });
}

export default function Navigation() {
  const { user, setUser } = useContext(UserContext) as IUserContext;

  if (isLoggedIn() && !loginAttemptMade) {
    getCurrentUser()
      .then((user) => {
        setUser(user);
      })
      .catch((err) => {
        if (!(err instanceof UserNotLoggedError)) {
          console.log("An error occurred while retrieving the current user. ", err)
        }
        redirect("/");
      }).finally(() => {
        loginAttemptMade = true;
      });
  }

    return (
        <>
            <nav className='navbar navbar-expand-lg navbar-light bg-light mb-3 p-3 justify-content-between'>
                <div className="container">
                    <Link to="/" className="navbar-brand">Movie Nights</Link>
                    <button className="navbar-toggler mr-3" type="button" data-bs-toggle="collapse" data-bs-target="#navbar" aria-controls="navbar" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbar">
                        <ul className="navbar-nav">
                            {isLoggedIn() &&
                                <li className="nav-item me-1">
                                    <a className="nav-link" data-active-if="/groups" href="/groups">Groups</a>
                                </li>
                            }
                            {isLoggedIn() &&
                                <li className="nav-item">
                                    <a className="nav-link" data-active-if="/movies" href="/movies">My Movies</a>
                                </li>
                            }
                        </ul>

                        <ul className="navbar-nav ms-auto">
                            { isLoggedIn() &&
                                <li className="nav-item p-0 m-0">
                                    <a className="nav-link p-0 m-0 active" href="#">
                                        <img 
                                            src="/api/avatar" 
                                            className="border rounded-circle me-1"
                                            style={{height: "45px", width: "45px", objectFit: "cover"}}
                                        />
                                        {user.name}
                                    </a>
                                </li>
                            }
                            { isLoggedIn() &&
                                <li className="nav-item">
                                    <a className="nav-link" href="/api/logout">Logout</a>
                                </li>
                            }
                            
                            { !isLoggedIn() &&
                                <li className="nav-item">
                                    <a className="nav-link" href="/">Login</a>
                                </li>
                            }
                        </ul>
                    </div>
                </div>
            </nav>
        </>
    )
}