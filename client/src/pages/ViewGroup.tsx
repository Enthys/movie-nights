import { Navigate, useParams } from "react-router-dom";
import IMovie from "../types/movie";
import {useEffect, useState} from "react";
import Movie from "../components/Movie";
import Group from "../types/group";

function getGroup(groupId: number): Promise<Group> {
    return new Promise((resolve, reject) => {
        fetch(`/api/groups/${groupId}`).then((resp) => {
            resp.json().then((data) => {
                if (resp.status !== 200) {
                    reject(new Error(data.error))
                    return
                }

                resolve(data.group)
            })
        })
    })
}

function getMovies(groupId: number, page: number): Promise<IMovie[]> {
    return new Promise((resolve, reject) => {
        fetch(`/api/groups/${groupId}/movies?page=${page}`).then((resp) => {
            resp.json().then((data) => {
                if (resp.status !== 200) {
                    reject(new Error(data.error))
                    return
                }

                resolve(data.movies)
            })
        })
    })
}

export default function ViewGroup() {
    const [error, setError] = useState("");
    const [group, setGroup] = useState({} as Group)
    const [page] = useState(1);
    const [movies, setMovies] = useState([] as IMovie[])
    const { groupId: groupIdStr } = useParams();
    const groupId = Number(groupIdStr);

    useEffect(() => {
        getGroup(groupId)
            .then((group) => {
                setGroup(group);
                getMovies(groupId, page)
                    .then(setMovies)
                    .catch(setError);
            })
            .catch(setError)
    }, [groupId, page])

    if (Number.isNaN(groupId)) {
        return <Navigate to="/groups" />;
    }

    

    return (<>
    <div className="container">
        {error !== "" &&
            <div className="text-danger">{error}</div>
        }
        <div className="row">
            <img src={group.name === '' ? '' : `https://api.dicebear.com/8.x/bottts/svg?seed=${group.name}`}
                className="col-sm-3 col-md-3 col-12"
                style={{maxHeight: '200px', objectFit: 'cover'}}
            />
            <div className="col-sm-9 col-md-9 col-12 ">
                    <h3 className="">{group.name}</h3>
                    <p>
                        {group.description}
                    </p>
                </div>
            </div>
            <hr />

            <div className="row mb-2">
                <div className="col-sm-12 col-lg-4">
                    <div className="input-group">
                        <input type="text" className="form-control" placeholder="Movie IMDb link" />
                        <button type="button" className="btn btn-outline-success">Add</button>
                    </div>
                </div>
                <div className="col-sm-12 col-lg-1">
                    <p className="my-1 lh-5 text-center"><span className="align-baseline">Or</span></p>
                </div>
                <div className="col-sm-12 col-lg-4">
                    <div className="input-group">
                        <input type="text" className="form-control" placeholder="Movie name" />
                        <button type="button" className="btn btn-outline-success">Seach</button>
                    </div>
                </div>
                <div className="col-sm-12 col-lg-1">
                    <p className="my-1 lh-5 text-center"><span className="align-baseline">Or</span></p>
                </div>
                <button className="px-2 px-sm-0 col-12 col-lg-2 text-center btn btn-success">Pick Random</button>
            </div>

            <div className="row">
              {movies.length > 0 && 
                movies.map((movie) => <Movie movie={movie}/>)
              }
            </div>
        </div>

    </>)
}
