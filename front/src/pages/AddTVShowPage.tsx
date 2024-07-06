import AddTVShowForm from '../components/AddTVShowForm';

const AddTVShowPage = () => {
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2 has-text-centered mb-6">Add New TV Show</h1>
        <AddTVShowForm />
      </div>
    </section>
  );
}

export default AddTVShowPage;
