import React, { useState, useEffect } from 'react';

const ReviewSection = ({ targetUUID }) => {
  const [reviews, setReviews] = useState([]);
  const [comment, setComment] = useState('');
  const [grade, setGrade] = useState(5); // Default grade value

  useEffect(() => {
    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    // Fetch reviews logic
    // This should interact with your backend to get reviews based on targetUUID
    console.log('Fetching reviews for', targetUUID);
    // Dummy data
    setReviews([
      { id: '1', userUUID: 'user1', grade: 5, comment: 'Great movie!' },
      {
        id: '2',
        userUUID: 'user2',
        grade: 4,
        comment: 'Enjoyable but a bit long.',
      },
    ]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    // Submit review logic
    console.log(`Submitting review: ${comment} with grade: ${grade}`);
    // Add logic to send this data to your backend
  };

  return (
    <div className="box">
      <h2 className="title is-4">Reviews</h2>
      {reviews.map((review) => (
        <div key={review.id} className="box">
          <p>
            <strong>User:</strong> {review.userUUID}
          </p>
          <p>
            <strong>Grade:</strong> {review.grade}
          </p>
          <p>{review.comment}</p>
        </div>
      ))}
      <div className="field">
        <label className="label">Your Review</label>
        <div className="control">
          <textarea
            className="textarea"
            placeholder="Add a review"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
          />
        </div>
      </div>
      <div className="field">
        <label className="label">Grade</label>
        <div className="control">
          <div className="select">
            <select value={grade} onChange={(e) => setGrade(e.target.value)}>
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
              <option value="4">4</option>
              <option value="5">5</option>
            </select>
          </div>
        </div>
      </div>
      <div className="field">
        <div className="control">
          <button
            type="button"
            className="button is-link"
            onClick={handleSubmit}
          >
            Submit Review
          </button>
        </div>
      </div>
    </div>
  );
};

export default ReviewSection;
