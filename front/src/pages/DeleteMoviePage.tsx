import MovieTable from '../components/MovieTable';

const DeleteMoviePage = () => {
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Delete Movie</h1>
        <MovieTable />
      </div>
    </section>
  );
}

export default DeleteMoviePage;
