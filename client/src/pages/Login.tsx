export default function Login() {
    return (
        <>
            <div className="container">
                <div className="row offset-lg-3 col-lg-6">
                    <div className="card">
                        <form method="GET" className="card-body" action="/api/login/google">
                            <div className="form-group mt-3">
                                <button type="submit" className="btn form-control btn-success">
                                    Sign in
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </>
    )
}