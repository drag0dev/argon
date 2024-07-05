import MovieTable from '../components/MovieTable';

const EditOrDeleteMoviePage = () => {
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Edit/Delete Movie</h1>
        <MovieTable />
      </div>
    </section>
  );
}

export default EditOrDeleteMoviePage;
