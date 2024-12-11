package ch.unisg.ics.interactions;

import cartago.OPERATION;
import cartago.ObsProperty;
import jason.asSyntax.ASSyntax;
import org.hyperagents.yggdrasil.cartago.artifacts.SyncedHypermediaTDArtifact;


public class SyncedSecuritySystem extends SyncedHypermediaTDArtifact {

    public void init() {
        super.init();
        defineObsProperty("isLocked", true);
        ObsProperty roomIsLocked = getObsProperty("isLocked");
        String artId = getArtifactId().toString();
        roomIsLocked.addAnnot(ASSyntax.createStructure("roomName", ASSyntax.createNumber(Double.parseDouble(artId.substring(artId.length() - 1)))));
    }

    public void init(Object initializationParameters) {
        super.init(initializationParameters);
        defineObsProperty("isLocked", true);
        ObsProperty roomIsLocked = getObsProperty("isLocked");
        String artId = getArtifactId().toString();
      roomIsLocked.addAnnot(ASSyntax.createStructure("roomName", ASSyntax.createNumber(Double.parseDouble(artId.substring(artId.length() - 1)))));
    }

    @OPERATION
    public void unlockRoom() {
        ObsProperty roomIsLocked = getObsProperty("isLocked");

        if (!roomIsLocked.booleanValue()) {
            log("Room is already unlocked");
            return;
        }

        updateInternalTimestamp(); //Conceptual point
        roomIsLocked.updateValue(false);

    }

    @OPERATION
    public void lockRoom() {
        ObsProperty roomIsLocked = getObsProperty("isLocked");

        if (roomIsLocked.booleanValue()) {
            log("Room is already locked");
            return;
        }

        updateInternalTimestamp();
        roomIsLocked.updateValue(true);

    }

    @Override
    protected void registerInteractionAffordances() {
      super.registerInteractionAffordances();
        // Register one action affordance with an input schema
        registerActionAffordance("http://example.org/unlockRoom", "unlockRoom", "unlockRoom");
        registerActionAffordance("http://example.org/lockRoom", "lockRoom", "lockRoom");
    }

}
