import "./App.css";
import { SignupFormDemo } from "./component/form";
import AuroraBackground from "./ui/aurira-background";
import { Navbar } from "./component/navbar";
import { Footer } from "./component/footer";
function App() {
  return (
    <>
    <Navbar />
      <AuroraBackground>
        <div className="grid h-screen w-screen">
          <SignupFormDemo />
        </div>
      </AuroraBackground>
      <Footer />
    </>
  );
}

export default App;
