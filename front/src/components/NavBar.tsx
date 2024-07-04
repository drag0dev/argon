import { Link } from 'react-router-dom';

const NavBar = () => {
  return (
    <nav className="navbar" aria-label="main navigation">
      <div className="navbar-brand">Joe</div>

      <a
        href="/settings"
        role="button"
        className="navbar-burger"
        aria-label="menu"
        aria-expanded="false"
        data-target="navbarBasicExample"
      >
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span style={{ display: 'none' }}>Menu</span>
      </a>

      <div id="navbarBasicExample" className="navbar-menu">
        <div className="navbar-start">
          <Link to="/home" className="navbar-item">
            Home
          </Link>
          <Link to="/docs" className="navbar-item">
            Documentation
          </Link>
          <div className="navbar-item has-dropdown is-hoverable">
            <Link to="/more" className="navbar-link">
              More
            </Link>
            <div className="navbar-dropdown">
              <Link to="/about" className="navbar-item">
                About
              </Link>
              <Link to="/joes" className="navbar-item">
                Joes
              </Link>
              <Link to="/contact" className="navbar-item">
                Contact
              </Link>
              <hr className="navbar-divider" />
              <Link to="report-issue" className="navbar-item">
                Report an issue
              </Link>
            </div>
          </div>
        </div>

        <div className="navbar-end">
          <div className="navbar-item">
            <div className="buttons">
              <Link to="sign-up" className="button is-primary">
                <strong>Sign up</strong>
              </Link>
              <Link to="log-in" className="button is-light">
                Log in
              </Link>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
