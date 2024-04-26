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
                        <h5 className="card-title">{movie.name}</h5>
                        <p className="card-text">
                            {movie.description}
                            Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the
                            industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type
                            and scrambled it to make a type specimen book. It has survive
                        </p>
                        <p>
                            {movie.genres.map((genre) => <span className="badge bg-secondary">{genre}</span>)}
                        </p>
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