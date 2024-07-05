import TVShowTable from '../components/TVShowTable';

const EditOrDeleteTVShowPage = () => {
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Edit/Delete TVShow</h1>
        <TVShowTable />
      </div>
    </section>
  );
}

export default EditOrDeleteTVShowPage;
