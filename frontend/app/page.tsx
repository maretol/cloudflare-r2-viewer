import { EventsEmit } from "@/wailsjs/runtime/runtime";
import Image from "next/image";

export default function Home() {
  const handleClick = () => {
    EventsEmit("testEvent", { name: "test" });
  };

  return (
    <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex"></div>
  );
}
