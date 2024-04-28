import {default as IMovie} from '../types/movie';

interface MovieParams {
    movie: IMovie;
}

export default function Movie({ movie }: MovieParams) {
    return (
        <>
            <div className="card g-0 mb-2 p-2">
                <div className="row justify-content-around">
                    <img
                        className="img-fluid rounded-start col-sm-2 col-5"
                        src={movie.avatar}
                        style={{objectFit: "contain", maxHeight: "250px"}}
                    />
                    <div className="col-sm-10 col-7">
                        <h5 className="card-title">{movie.name} <span style={{color: "#a38d00"}}>({movie.imdbRating}/10)</span></h5>

                        <p className="card-text">
                            {movie.description}
                        </p>

                        <div className="d-flex justify-content-start mb-2">
                            {movie.genres.map((genre) => (<div className="badge bg-secondary me-1">{genre}</div>))}
                        </div>

                        <div className="row justify-content-around">
                            <button type="button" className="col-5 col-sm-3 col-md-4 col-lg-3 btn btn-primary">IMDb</button>
                            <button type="button" className="col-5 col-sm-3 col-md-4 ms-1 me-3 btn btn-success">Select</button>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}