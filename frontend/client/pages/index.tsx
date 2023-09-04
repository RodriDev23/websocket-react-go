import ChatForm from './components/chatForm'

export default function Home() {
  return (
   <div className='flex  justify-center items-center gap-10 min-w-screen min-h-screen'>
   <ChatForm connection='server1' nameServer="chat 1" />
   <ChatForm connection='server2' nameServer="chat 2" />
   <ChatForm connection='server3' nameServer="chat 3" />
   </div>
  )
}
