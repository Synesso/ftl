import { Timestamp } from '@bufbuild/protobuf'
import { useContext, useEffect, useState } from 'react'
import { AttributeBadge } from '../../../components/AttributeBadge'
import { CloseButton } from '../../../components/CloseButton'
import { CodeBlock } from '../../../components/CodeBlock'
import { useClient } from '../../../hooks/use-client'
import { ConsoleService } from '../../../protos/xyz/block/ftl/v1/console/console_connect'
import { CallEvent } from '../../../protos/xyz/block/ftl/v1/console/console_pb'
import { SidePanelContext } from '../../../providers/side-panel-provider'
import { getRequestCalls } from '../../../services/console.service'
import { formatDuration } from '../../../utils/date.utils'
import { DeploymentCard } from '../../deployments/DeploymentCard'
import { RequestGraph } from '../../requests/RequestGraph'
import { verbRefString } from '../../verbs/verb.utils'
import { TimelineTimestamp } from './TimelineTimestamp'

export const TimelineCallDetails = ({ timestamp, call }: { timestamp: Timestamp; call: CallEvent }) => {
  const client = useClient(ConsoleService)
  const { closePanel } = useContext(SidePanelContext)
  const [requestCalls, setRequestCalls] = useState<CallEvent[]>([])
  const [selectedCall, setSelectedCall] = useState(call)

  useEffect(() => {
    setSelectedCall(call)
  }, [call])

  useEffect(() => {
    const abortController = new AbortController()
    const fetchRequestCalls = async () => {
      if (selectedCall.requestName === undefined) {
        return
      }
      const calls = await getRequestCalls({
        abortControllerSignal: abortController.signal,
        requestKey: selectedCall.requestName,
      })
      setRequestCalls(calls.reverse())
    }

    fetchRequestCalls()

    return () => {
      abortController.abort()
    }
  }, [client, selectedCall])

  return (
    <div className='p-4'>
      <div className='flex items-center justify-between'>
        <div className='flex items-center space-x-2'>
          <div className=''>
            {call.destinationVerbRef && (
              <div
                className={`inline-block rounded-md bg-indigo-200 dark:bg-indigo-700 px-2 py-1 mr-1 text-sm font-medium text-gray-700 dark:text-gray-100`}
              >
                {verbRefString(call.destinationVerbRef)}
              </div>
            )}
          </div>
          <TimelineTimestamp timestamp={timestamp} />
        </div>
        <CloseButton onClick={closePanel} />
      </div>

      <div className='pt-4'>
        <RequestGraph calls={requestCalls} call={selectedCall} setSelectedCall={setSelectedCall} />
      </div>

      <div className='text-sm pt-2'>Request</div>
      <CodeBlock code={JSON.stringify(JSON.parse(selectedCall.request), null, 2)} language='json' />

      {selectedCall.response != 'null' && (
        <>
          <div className='text-sm pt-2'>Response</div>
          <CodeBlock code={JSON.stringify(JSON.parse(selectedCall.response), null, 2)} language='json' />
        </>
      )}

      {selectedCall.error && (
        <>
          <h3 className='pt-4'>Error</h3>
          <CodeBlock code={selectedCall.error} language='text' />
          {selectedCall.stack && (
            <>
              <h3 className='pt-4'>Stack</h3>
              <CodeBlock code={selectedCall.stack} language='text' />
            </>
          )}
        </>
      )}

      <DeploymentCard className='mt-4' deploymentName={call.deploymentName} />

      <ul className='pt-4 space-y-2'>
        {selectedCall.requestName && (
          <li>
            <AttributeBadge name='Request' value={selectedCall.requestName} />
          </li>
        )}
        <li>
          <AttributeBadge name='Duration' value={formatDuration(selectedCall.duration)} />
        </li>
        {selectedCall.destinationVerbRef && (
          <li>
            <AttributeBadge name='Destination' value={verbRefString(selectedCall.destinationVerbRef)} />
          </li>
        )}
        {selectedCall.sourceVerbRef && (
          <li>
            <AttributeBadge name='Source' value={verbRefString(selectedCall.sourceVerbRef)} />
          </li>
        )}
      </ul>
    </div>
  )
}
