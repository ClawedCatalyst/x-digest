export default function Home() {
  return (
    <main>
      <h1>X Daily Digest MVP</h1>
      <p>Login with X to generate your daily digest.</p>

      <a
        href="http://localhost:8080/auth/x/start"
        style={{
          display: "inline-block",
          padding: "10px 14px",
          border: "1px solid #333",
          borderRadius: 10,
          textDecoration: "none",
        }}
      >
        Login with X
      </a>
    </main>
  );
}